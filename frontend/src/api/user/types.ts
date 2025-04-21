export interface IUserResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  birthDate: string;
  createdAt: Date;
  updatedAt: Date;
  friendCount: number;
}

export interface IUserProfileResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  birthDate: string;
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
  createdAt: Date;
  updatedAt: Date;
}

export interface IFriendsResponse {
  friends: {
    id: string;
    firstName: string;
    birthDate: Date;
    profilePicture: string;
    friendCount: number;
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

export interface ICommunityResponse {
  name: string;
  memberCount: number;
}

export interface ICommunityProfileResponse {
  description: string;
  profilePicture: string;
}

export interface ICommunityMember {
  id: number;
  firstName: string;
  profilePicture: string;
  friendCount: number;
}

export interface ICommunityMembersResponse {
  members: ICommunityMember[];
  memberCount: number;
}

export interface IRelatedCommunitiesResponse {
  communities: any[];
  memberCount: number;
}

export interface ICreateScrapRequest {
  userAuthorId: number;
  userRecipientId: number;
  content: string;
}

export interface IScrapsResponse {
  scraps: {
    id: number;
    content: string;
  }[];
}
