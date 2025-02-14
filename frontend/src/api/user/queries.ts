import arrasa from "assets/communities/arrasa.png";
import borracha from "assets/communities/borracha.png";
import campainha from "assets/communities/campainha.png";
import cinco from "assets/communities/cinco.png";
import duplas from "assets/communities/duplas.png";
import gina from "assets/communities/gina.png";
import odeio_acordar_cedo from "assets/communities/odeio_acordar_cedo.png";
import sorvete from "assets/communities/sorvete.png";
import tony from "assets/communities/tony.png";
import aline from "assets/users/aline.jpg";
import caique from "assets/users/caique.jpg";
import carol from "assets/users/carol.jpg";
import diego from "assets/users/diego.jpg";
import frizz from "assets/users/frizz.jpg";
import ligia from "assets/users/ligia.jpg";
import luciana from "assets/users/luciana.jpg";
import renato from "assets/users/renato.jpg";
import weslley from "assets/users/weslley.jpg";

export const useUserDetails = (userId: number) => ({
  userId: userId,
  firstName: "Frizz",
  lastName: "Chapper",
  profilePictureUrl: frizz,
  accountCreatedAt: new Date(),
});

export const useUserStatistics = (userId: number) => ({
  userId: userId,
  totalProfileViews: 3558,
  lastWeekProfileViews: 45,
  totalFriends: 111,
  totalCommunities: 123,
  undirectConnections: 34537335,
  totalScraps: 123,
  totalPhotos: 12,
  totalVideos: 1,
  totalFans: 8,
  trustworthy: 1,
  cool: 1,
  beautiful: 1,
  profileLastVisitors: [
    { name: "Aline", birthDay: new Date(), profilePicture: aline, friends: 123 },
    { name: "Caique", birthDay: new Date(), profilePicture: caique, friends: 123 },
    { name: "Carol", birthDay: new Date(), profilePicture: carol, friends: 123 },
    { name: "Diego", birthDay: new Date(), profilePicture: diego, friends: 123 },
    { name: "Lígia", birthDay: new Date(), profilePicture: ligia, friends: 123 },
    { name: "Luciana", birthDay: new Date(), profilePicture: luciana, friends: 123 },
    { name: "Renato", birthDay: new Date(), profilePicture: renato, friends: 123 },
    { name: "Weslley", birthDay: new Date(), profilePicture: weslley, friends: 123 },
  ],
});

export const useUserFriends = (userId: number) => [
  { id: 1, name: "Aline", birthDay: new Date(), profilePicture: aline, friends: 123 },
  { id: 2, name: "Caique", birthDay: new Date(), profilePicture: caique, friends: 123 },
  { id: 3, name: "Carol", birthDay: new Date(), profilePicture: carol, friends: 123 },
  { id: 4, name: "Diego", birthDay: new Date(), profilePicture: diego, friends: 123 },
  { id: 5, name: "Lígia", birthDay: new Date(), profilePicture: ligia, friends: 123 },
  { id: 6, name: "Luciana", birthDay: new Date(), profilePicture: luciana, friends: 123 },
  { id: 7, name: "Renato", birthDay: new Date(), profilePicture: renato, friends: 123 },
  { id: 8, name: "Weslley", birthDay: new Date(), profilePicture: weslley, friends: 123 },
];

export const useUserCommunities = (userId: number) => [
  { id: 1, name: "Deus me disse: desce e arrasa", picture: arrasa, members: 450 },
  { id: 2, name: "Eu nunca terminei uma borracha", picture: borracha, members: 3214 },
  { id: 3, name: "Tocava campainha e corria", picture: campainha, members: 54352 },
  { id: 4, name: "Só mais 5 minutinhos", picture: cinco, members: 123 },
  { id: 5, name: "Não sei individualizar duplas", picture: duplas, members: 12 },
  { id: 6, name: "Tive medo da Gina dos palitos", picture: gina, members: 9 },
  { id: 7, name: "Eu odeio acordar cedo", picture: odeio_acordar_cedo, members: 839129 },
  { id: 8, name: "Queria sorvete, mas era feijão", picture: sorvete, members: 7 },
  { id: 9, name: "Eu colho flores e o Tony Ramos", picture: tony, members: 455 },
];

export const useTodaysLuck = (userId: number) => "Boas notícias chegarão por e-mail";
