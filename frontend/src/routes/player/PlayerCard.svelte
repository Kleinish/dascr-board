<script>
  export let uid, name, nickname, image, onDelete, showDelete, active;
  import myenv from '$utils/env';
  import { onMount } from 'svelte';
  import api from '$utils/api';

  let apiBase = myenv.apiBase;
  let activeStyle =
    'border-4 border-white shadow-2xl opacity-100 transition duration-700 ease-in-out';
  
  let stats = {
    gamesPlayed: 0,
    averageScore: 0,
    highestScore: 0,
    totalThrows: 0,
    averageThrow: 0
  };

  onMount(async () => {
    try {
      // Fetch games to calculate statistics
      const res = await api.get('game');
      const games = await res.json();
      
      let totalScore = 0;
      let highestScore = 0;
      let totalThrows = 0;
      let gamesPlayed = 0;

      games.forEach(game => {
        if (game.Players && game.Players.includes(uid)) {
          gamesPlayed++;
          const playerStats = game.Statistics?.[uid];
          if (playerStats) {
            totalScore += playerStats.Score || 0;
            totalThrows += playerStats.TotalThrowCount || 0;
            if (playerStats.Score > highestScore) {
              highestScore = playerStats.Score;
            }
          }
        }
      });

      stats = {
        gamesPlayed,
        averageScore: gamesPlayed ? Math.round(totalScore / gamesPlayed) : 0,
        highestScore,
        totalThrows,
        averageThrow: totalThrows ? Math.round((totalScore / totalThrows) * 10) / 10 : 0
      };
    } catch (error) {
      console.error('Error fetching player statistics:', error);
    }
  });
</script>

<div
  class="flex flex-col lg:flex-row bg-black bg-opacity-50 rounded-lg overflow-hidden shadow-lg opacity-80 {active && activeStyle} playercard">
  <slot name="avatar">
    <div class="flex-none xs: hidden lg:block">
      <img class="" width="300" height="300" src={apiBase + image} alt="" />
    </div>
  </slot>
  <div class="m-6 flex-none w-60 overflow-hidden">
    <p class="font-semibold text-4xl">{name}</p>
    {#if nickname}
      <p class="text-gray-400 text-3xl">{nickname}</p>
    {:else}<br />{/if}
    {#if $$slots.points}
      <slot name="points" />
    {/if}

    {#if showDelete}
      <button
        class="block uppercase border-2 hover:bg-black hover:bg-opacity-30 text-lg p-4 mt-4 rounded-2xl"
        on:click={onDelete(uid)}><i class="fas fa-trash"></i>
        Delete</button>
    {/if}
  </div>
  <div class="m-6 w-full">
    {#if $$slots.score}
      <slot name="score" />
    {:else}
      <div class="grid grid-cols-2 gap-4 text-lg">
        <div class="stats-item">
          <span class="text-gray-400">Games Played:</span>
          <span class="font-semibold">{stats.gamesPlayed}</span>
        </div>
        <div class="stats-item">
          <span class="text-gray-400">Average Score:</span>
          <span class="font-semibold">{stats.averageScore}</span>
        </div>
        <div class="stats-item">
          <span class="text-gray-400">Highest Score:</span>
          <span class="font-semibold">{stats.highestScore}</span>
        </div>
        <div class="stats-item">
          <span class="text-gray-400">Total Throws:</span>
          <span class="font-semibold">{stats.totalThrows}</span>
        </div>
        <div class="stats-item col-span-2">
          <span class="text-gray-400">Average Per Throw:</span>
          <span class="font-semibold">{stats.averageThrow}</span>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .stats-item {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem;
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 0.5rem;
  }
</style>
