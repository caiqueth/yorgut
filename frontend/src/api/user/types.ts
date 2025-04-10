export interface IUserResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  birthDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface IFriendsResponse {
  friends: {
    id: string;
    firstName: string;
    birthDate: Date;
    profilePicture: string;
    frienshipsCount?: number;
  }[];
}

export interface ICommunityMembershipsResponse {
  communities: {
    id: number;
    name: string;
    profilePicture: string;
    members?: number;
  }[];
}