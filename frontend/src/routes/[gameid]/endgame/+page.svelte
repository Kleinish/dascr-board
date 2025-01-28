<script>
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import api from '../../../utils/api';
    import myenv from '../../../utils/env';

    let gameid = $page.params.gameid;
    let players = [];
    let gameData = {};
    let apiBaseURL = myenv.apiBase;

    onMount(async () => {
        try {
            const gameRes = await api.get(`game/${gameid}`);
            gameData = await gameRes.json();

            const playersRes = await api.get(`game/${gameid}/players`);
            players = await playersRes.json();

            // Sort players based on game type
            switch (gameData.Game) {
                case 'x01':
                case 'elim':
                    players.sort((a, b) => a.Score - b.Score);
                    break;
                case 'cricket':
                    players.sort((a, b) => b.Score - a.Score);
                    break;
                case 'atc':
                case 'split':
                case 'shanghai':
                    players.sort((a, b) => b.Score - a.Score);
                    break;
                default:
                    break;
            }
        } catch (error) {
            console.error('Error fetching game data:', error);
        }
    });

    function getPlayerRank(index) {
        const ranks = ['🥇 1st', '🥈 2nd', '🥉 3rd'];
        return index < 3 ? ranks[index] : `${index + 1}th`;
    }
</script>

<div class="container mx-auto p-4 text-white">
    <h1 class="text-3xl font-bold mb-6 text-center">
        Game Over - Final Leaderboard
    </h1>

    <div class="grid grid-cols-1 gap-4 max-w-2xl mx-auto">
        {#each players as player, index}
            <div 
                class="bg-black bg-opacity-50 rounded-lg p-4 flex items-center space-x-4 
                       {index === 0 ? 'border-4 border-yellow-500' : 
                        index === 1 ? 'border-4 border-gray-400' : 
                        index === 2 ? 'border-4 border-yellow-700' : ''}"
            >
                <div class="flex-shrink-0">
                    <img 
                        src="{apiBaseURL}{player.Image}" 
                        alt="{player.Name}" 
                        class="w-16 h-16 rounded-full object-cover"
                    >
                </div>
                <div class="flex-grow">
                    <div class="flex justify-between items-center">
                        <div>
                            <h2 class="text-xl font-bold">
                                {getPlayerRank(index)} - {player.Name} 
                                {player.Nickname ? `(${player.Nickname})` : ''}
                            </h2>
                            <p class="text-sm opacity-75">
                                {#if gameData.Game === 'x01' || gameData.Game === 'elim'}
                                    Score: {player.Score} | Average: {player.Average}
                                {:else if gameData.Game === 'cricket'}
                                    Points Closed: {player.Score}
                                {:else}
                                    Score: {player.Score}
                                {/if}
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        {/each}
    </div>

    <div class="text-center mt-8">
        <a 
            href="/" 
            class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded inline-block"
        >
            Return to Home
        </a>
    </div>
</div>