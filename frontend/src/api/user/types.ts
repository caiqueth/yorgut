export interface IUserResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  birthDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface IUserProfileResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  birthDate: string;
  createdAt: string;
  updatedAt: string;
  profilePicture?: string;
  description?: string;
  relationshipStatus?: string;
  interests?: string;
  about?: string;
  religion?: string;
  humor?: string;
  fashion?: string;
  livesWith?: string;
  birthCity?: string;
  website?: string;
  passions?: string;
  sports?: string;
  hobbies?: string;
  books?: string;
  music?: string;
  movies?: string;
  cuisine?: string;
}

export interface IFriendsResponse {
  friends: {
    id: string;
    firstName: string;
    birthDate: Date;
    profilePicture: string;
    frienshipCount?: number;
  }[];
}

export interface IPendingFriendshipsResponse {
  pendingFriendships: {
    id: string;
    firstName: string;
    lastName: string;
    profilePicture: string;
  }[];
}

export interface ICommunityMembershipsResponse {
  communities: {
    id: number;
    name: string;
    profilePicture: string;
    membershipCount: number;
  }[];
}
