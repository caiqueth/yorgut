// "1 amigo", "3 amigos".
export const friendCountText = (count: number) => `${count} ${count === 1 ? "amigo" : "amigos"}`;

// "1 membro", "3 membros".
export const memberCountText = (count: number) => `${count} ${count === 1 ? "membro" : "membros"}`;
