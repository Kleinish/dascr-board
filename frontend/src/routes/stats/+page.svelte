<script>
  import api from '$utils/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let games = [];
  let players = {};
  let loading = true;

  onMount(async () => {
    try {
      // Fetch all games and players
      const [gamesRes, playersRes] = await Promise.all([
        api.get('game'),
        api.get('player')
      ]);
      
      const gamesData = await gamesRes.json();
      const playersData = await playersRes.json();

      // Create a map of player IDs to player data
      playersData.forEach(player => {
        players[player.UID] = player;
      });

      // Process games to include player names and format dates
      games = gamesData
        .filter(game => game.EndTime) // Only show finished games
        .map(game => ({
          ...game,
          formattedDate: new Date(game.StartTime).toLocaleDateString(),
          formattedTime: new Date(game.StartTime).toLocaleTimeString(),
          duration: calculateDuration(game.StartTime, game.EndTime),
          playerNames: game.Players.map(pid => players[pid]?.Name || 'Unknown Player'),
          winner: game.Winner ? players[game.Winner]?.Name || 'Unknown Player' : 'No Winner'
        }))
        .sort((a, b) => new Date(b.EndTime) - new Date(a.EndTime));
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      loading = false;
    }
  });

  function calculateDuration(start, end) {
    const duration = new Date(end) - new Date(start);
    const minutes = Math.floor(duration / 60000);
    const seconds = Math.floor((duration % 60000) / 1000);
    return `${minutes}m ${seconds}s`;
  }

  function viewGameStats(gameId) {
    goto(`/stats/${gameId}`);
  }
</script>

<div class="max-w-7xl mx-auto px-4 py-8">
  <h1 class="text-4xl font-bold mb-8 text-center">Game History</h1>

  {#if loading}
    <div class="text-center text-xl">Loading...</div>
  {:else if games.length === 0}
    <div class="text-center text-xl">No completed games found.</div>
  {:else}
    <div class="grid gap-6">
      {#each games as game}
        <div class="bg-black bg-opacity-50 rounded-lg p-6 shadow-lg">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div>
              <h3 class="text-xl font-semibold">{game.Type}</h3>
              <p class="text-gray-400">{game.formattedDate} at {game.formattedTime}</p>
            </div>
            <div>
              <p class="text-gray-400">Duration</p>
              <p class="font-semibold">{game.duration}</p>
            </div>
            <div>
              <p class="text-gray-400">Players</p>
              <p class="font-semibold">{game.playerNames.join(', ')}</p>
            </div>
            <div>
              <p class="text-gray-400">Winner</p>
              <p class="font-semibold">{game.winner}</p>
            </div>
          </div>
          <div class="mt-4 flex justify-end">
            <button
              class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
              on:click={() => viewGameStats(game.ID)}>
              View Details
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>