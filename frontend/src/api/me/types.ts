export interface IMeResponse {
  birthDate: string;
  createdAt: string;
  email: string;
  firstName: string;
  friendCount: number;
  id: number;
  lastName: string;
  updatedAt: string;
}

export interface IMeUpdateResponse {
  data: string;
}

export interface IMeFriendsResponse {
  friends: {
    birthDate: string;
    firstName: string;
    friendCount: number;
    id: number;
    profilePicture: string;
  }[];
}

export interface IMeFriendsPendingResponse {
  friends: {
    birthDate: "string";
    firstName: "string";
    friendCount: 0;
    id: 0;
    profilePicture: "string";
  }[];
}

export interface IJoinCommunityResponse {
  data: string;
}

export interface ILeaveCommunityResponse {
  data: string;
}

export interface IMeMembershipsResponse {
  communities: {
    firstName: number;
    id: number;
    membershipCount: number;
    name: string;
    profilePicture: string;
  }[];
}

export interface IMePhotoAlbumsResponse {
  photoAlbums: {
    name: string;
    userId: number;
  }[];
}

export interface IMeProfileResponse {
  about: string;
  birthCity: string;
  birthDate: string;
  books: string;
  createdAt: string;
  cuisine: string;
  description: string;
  email: string;
  fashion: string;
  firstName: string;
  friendCount: number;
  hobbies: string;
  humor: string;
  id: number;
  interests: string;
  lastName: string;
  livesWith: string;
  movies: string;
  music: string;
  passions: string;
  profilePicture: string;
  relationshipStatus: string;
  religion: string;
  sports: string;
  updatedAt: string;
  website: string;
}
