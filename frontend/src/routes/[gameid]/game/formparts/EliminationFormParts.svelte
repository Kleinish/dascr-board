<script>
  import { onMount } from 'svelte';
  import setupGame from '$stores/gameStore';

  onMount(() => {
    $setupGame.variant = '301';
  });
</script>

<style>
    /* Removed unused CSS selectors */
</style>

<div class="flex flex-col">
  <label for="selectVariant" class="uppercase font-bold text-lg"
    >Which Elimination do you want to play:</label
  >
  <select
    bind:value={$setupGame.variant}
    class="border py-2 px-3 text-gray-900"
    id="selectVariant"
    name="variant"
  >
    <option value="301" selected="selected">301</option>
    <option value="501">501</option>
    <option value="701">701</option>
    <option value="901">901</option>
  </select>
</div>
<div class="flex flex-col">
  <label for="selectInCondition" class="uppercase font-bold text-lg"
    >Condition to get into the game:</label
  >
  <select
    bind:value={$setupGame.in}
    class="border py-2 px-3 text-gray-900"
    id="selectInCondition"
    name="inCondition"
  >
    <option value="straight">Straight</option>
    <option value="double">Double</option>
    <option value="master">Master</option>
  </select>
</div>
<div class="flex flex-col">
  <label for="selectOutCondition" class="uppercase font-bold text-lg"
    >Condition to end the game:</label
  >
  <select
    bind:value={$setupGame.out}
    class="border py-2 px-3 text-gray-900"
    id="selectOutCondition"
    name="outCondition"
  >
    <option value="straight">Straight</option>
    <option value="double" selected="selected">Double</option>
    <option value="master">Master</option>
  </select>
</div>
