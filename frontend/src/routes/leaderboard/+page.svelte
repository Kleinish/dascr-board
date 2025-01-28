<script>
  import api from '$utils/api';
  import { onMount } from 'svelte';

  let players = [];
  let games = [];
  let loading = true;
  let selectedGameType = 'all';
  let timeRange = 'all';

  let playerStats = new Map();

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    try {
      // Fetch all games and players
      const [gamesRes, playersRes] = await Promise.all([
        api.get('game'),
        api.get('player')
      ]);
      
      games = await gamesRes.json();
      players = await playersRes.json();

      calculateStats();
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      loading = false;
    }
  }

  function calculateStats() {
    playerStats.clear();

    // Initialize stats for all players
    players.forEach(player => {
      playerStats.set(player.UID, {
        name: player.Name,
        nickname: player.Nickname,
        gamesPlayed: 0,
        gamesWon: 0,
        totalScore: 0,
        highestScore: 0,
        averageScore: 0,
        totalThrows: 0,
        averagePerThrow: 0,
        winRate: 0
      });
    });

    // Filter games based on selected type and time range
    const filteredGames = games.filter(game => {
      if (!game.EndTime) return false; // Only include finished games
      
      if (selectedGameType !== 'all' && game.Type !== selectedGameType) return false;

      if (timeRange !== 'all') {
        const gameDate = new Date(game.EndTime);
        const now = new Date();
        const thirtyDaysAgo = new Date(now.setDate(now.getDate() - 30));
        const ninetyDaysAgo = new Date(now.setDate(now.getDate() - 90));
        
        if (timeRange === '30days' && gameDate < thirtyDaysAgo) return false;
        if (timeRange === '90days' && gameDate < ninetyDaysAgo) return false;
      }

      return true;
    });

    // Calculate stats for each player
    filteredGames.forEach(game => {
      game.Players.forEach(playerId => {
        const stats = playerStats.get(playerId);
        if (!stats) return;

        stats.gamesPlayed++;
        if (game.Winner === playerId) stats.gamesWon++;

        const playerGameStats = game.Statistics?.[playerId];
        if (playerGameStats) {
          stats.totalScore += playerGameStats.Score || 0;
          stats.totalThrows += playerGameStats.TotalThrowCount || 0;
          if ((playerGameStats.Score || 0) > stats.highestScore) {
            stats.highestScore = playerGameStats.Score;
          }
        }
      });
    });

    // Calculate averages and win rates
    playerStats.forEach(stats => {
      if (stats.gamesPlayed > 0) {
        stats.averageScore = Math.round(stats.totalScore / stats.gamesPlayed);
        stats.winRate = Math.round((stats.gamesWon / stats.gamesPlayed) * 100);
      }
      if (stats.totalThrows > 0) {
        stats.averagePerThrow = Math.round((stats.totalScore / stats.totalThrows) * 10) / 10;
      }
    });

    // Convert Map to array and sort by win rate
    playerStats = new Map([...playerStats.entries()].sort((a, b) => b[1].winRate - a[1].winRate));
  }

  function handleFilterChange() {
    calculateStats();
  }

  $: {
    selectedGameType;
    timeRange;
    handleFilterChange();
  }
</script>

<div class="max-w-7xl mx-auto px-4 py-8">
  <h1 class="text-4xl font-bold mb-8 text-center">Leaderboard</h1>

  {#if loading}
    <div class="text-center text-xl">Loading...</div>
  {:else}
    <div class="mb-8 flex flex-wrap gap-4 justify-center">
      <div class="flex items-center gap-2">
        <label for="gameType">Game Type:</label>
        <select
          id="gameType"
          class="bg-black bg-opacity-50 rounded-lg px-4 py-2"
          bind:value={selectedGameType}>
          <option value="all">All Games</option>
          <option value="x01">X01</option>
          <option value="cricket">Cricket</option>
          <option value="atc">Around the Clock</option>
          <option value="shanghai">Shanghai</option>
        </select>
      </div>

      <div class="flex items-center gap-2">
        <label for="timeRange">Time Range:</label>
        <select
          id="timeRange"
          class="bg-black bg-opacity-50 rounded-lg px-4 py-2"
          bind:value={timeRange}>
          <option value="all">All Time</option>
          <option value="30days">Last 30 Days</option>
          <option value="90days">Last 90 Days</option>
        </select>
      </div>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="bg-black bg-opacity-70">
            <th class="px-4 py-2 text-left">Rank</th>
            <th class="px-4 py-2 text-left">Player</th>
            <th class="px-4 py-2 text-right">Games</th>
            <th class="px-4 py-2 text-right">Win Rate</th>
            <th class="px-4 py-2 text-right">Avg Score</th>
            <th class="px-4 py-2 text-right">Highest Score</th>
            <th class="px-4 py-2 text-right">Avg/Throw</th>
          </tr>
        </thead>
        <tbody>
          {#each [...playerStats.entries()] as [id, stats], index}
            {#if stats.gamesPlayed > 0}
              <tr class="border-b border-gray-800 hover:bg-black hover:bg-opacity-50">
                <td class="px-4 py-2">{index + 1}</td>
                <td class="px-4 py-2">
                  {stats.name}
                  {#if stats.nickname}
                    <span class="text-gray-400">({stats.nickname})</span>
                  {/if}
                </td>
                <td class="px-4 py-2 text-right">{stats.gamesPlayed}</td>
                <td class="px-4 py-2 text-right">{stats.winRate}%</td>
                <td class="px-4 py-2 text-right">{stats.averageScore}</td>
                <td class="px-4 py-2 text-right">{stats.highestScore}</td>
                <td class="px-4 py-2 text-right">{stats.averagePerThrow}</td>
              </tr>
            {/if}
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>