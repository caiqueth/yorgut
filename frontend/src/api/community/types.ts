export interface ICommunityResponse {
  createdBy: number;
  name: string;
}

export interface ICommunitiesResponse {
  communities: {
    firstName: number;
    id: number;
    membershipCount: number;
    name: string;
    owner: number;
    profilePicture: string;
  }[];
}

export interface ICommunityMembersResponse {
  memberCount: number;
  members: {
    firstName: string;
    friendCount: number;
    id: number;
    profilePicture: string;
  }[];
}

export interface ICommunityProfileResponse {
  description: string;
  profilePicture: string;
}

export interface ICommunityRelatedResponse {
  communities: {
    firstName: number;
    id: number;
    membershipCount: number;
    name: string;
    owner: number;
    profilePicture: string;
  }[];
}
