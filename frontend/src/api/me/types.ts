export interface IMeResponse {
  id: number;
  birthDate: Date;
  createdAt: Date;
  email: string;
  firstName: string;
  friendCount: number;
  lastName: string;
  updatedAt: string;
}

export interface IMeUpdateResponse {
  data: string;
}

export interface IMeFriendsResponse {
  friends: {
    birthDate: Date;
    firstName: string;
    friendCount: number;
    id: number;
    profilePicture: string;
  }[];
}

export interface IMeFriendsPendingResponse {
  count: number;
  pendingFriendships: {
    id: number;
    userRequestedId: number;
    firstName: string;
    lastName: string;
    profilePicture: string;
  }[];
}

export interface IJoinCommunityResponse {
  data: string;
}

export interface ILeaveCommunityResponse {
  data: string;
}

export interface IMeMembershipsResponse {
  membershipCount: number;
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

export interface IAmIMemberResponse {
  amIMember: boolean;
}

export interface IAmIFriendResponse {
  friendshipId?: number;
  amIFriend: boolean;
  isPendingMyReply: boolean;
  isPendingTheirReply: boolean;
}
