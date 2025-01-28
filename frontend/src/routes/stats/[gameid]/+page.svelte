<script>
  import api from '$utils/api';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';

  let game = null;
  let players = {};
  let loading = true;

  onMount(async () => {
    try {
      // Fetch game and players data
      const [gameRes, playersRes] = await Promise.all([
        api.get(`game/${$page.params.gameid}`),
        api.get('player')
      ]);
      
      game = await gameRes.json();
      const playersData = await playersRes.json();

      // Create a map of player IDs to player data
      playersData.forEach(player => {
        players[player.UID] = player;
      });

      // Add formatted dates
      game.formattedStartTime = new Date(game.StartTime).toLocaleString();
      game.formattedEndTime = new Date(game.EndTime).toLocaleString();
      game.duration = calculateDuration(game.StartTime, game.EndTime);
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
</script>

<div class="max-w-7xl mx-auto px-4 py-8">
  <h1 class="text-4xl font-bold mb-8 text-center">Game Statistics</h1>

  {#if loading}
    <div class="text-center text-xl">Loading...</div>
  {:else if !game}
    <div class="text-center text-xl">Game not found.</div>
  {:else}
    <div class="bg-black bg-opacity-50 rounded-lg p-6 shadow-lg mb-8">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div>
          <h3 class="text-xl font-semibold">{game.Type}</h3>
          <p class="text-gray-400">Game Type</p>
        </div>
        <div>
          <p class="font-semibold">{game.formattedStartTime}</p>
          <p class="text-gray-400">Start Time</p>
        </div>
        <div>
          <p class="font-semibold">{game.formattedEndTime}</p>
          <p class="text-gray-400">End Time</p>
        </div>
        <div>
          <p class="font-semibold">{game.duration}</p>
          <p class="text-gray-400">Duration</p>
        </div>
      </div>

      <div class="grid gap-6">
        <h2 class="text-2xl font-semibold">Player Statistics</h2>
        {#each game.Players as playerId}
          {@const player = players[playerId]}
          {@const stats = game.Statistics[playerId]}
          {#if player && stats}
            <div class="bg-black bg-opacity-30 rounded-lg p-4">
              <div class="flex items-center gap-4 mb-4">
                <h3 class="text-xl font-semibold">
                  {player.Name}
                  {#if player.Nickname}
                    <span class="text-gray-400">({player.Nickname})</span>
                  {/if}
                </h3>
                {#if game.Winner === playerId}
                  <span class="bg-green-600 text-white px-2 py-1 rounded text-sm">Winner</span>
                {/if}
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <div>
                  <p class="font-semibold">{stats.Score}</p>
                  <p class="text-gray-400">Final Score</p>
                </div>
                <div>
                  <p class="font-semibold">{stats.TotalThrowCount}</p>
                  <p class="text-gray-400">Total Throws</p>
                </div>
                <div>
                  <p class="font-semibold">{stats.TotalThrowCount ? Math.round((stats.Score / stats.TotalThrowCount) * 10) / 10 : 0}</p>
                  <p class="text-gray-400">Average Per Throw</p>
                </div>
                <div>
                  <p class="font-semibold">{stats.ThrowSum}</p>
                  <p class="text-gray-400">Total Points</p>
                </div>
              </div>

              {#if stats.ThrowRounds && stats.ThrowRounds.length > 0}
                <div class="mt-4">
                  <h4 class="font-semibold mb-2">Throw History</h4>
                  <div class="grid grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-2">
                    {#each stats.ThrowRounds as round, index}
                      <div class="bg-black bg-opacity-20 p-2 rounded text-center">
                        <p class="text-sm text-gray-400">Round {index + 1}</p>
                        <p class="font-semibold">{round.Throws.map(t => t.Value).join(', ')}</p>
                      </div>
                    {/each}
                  </div>
                </div>
              {/if}
            </div>
          {/if}
        {/each}
      </div>
    </div>
  {/if}
</div>