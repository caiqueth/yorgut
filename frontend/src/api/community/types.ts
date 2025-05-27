export interface ICommunityResponse {
  name: string;
  createdById: number;
  ownerId: number;
  ownerFirstName: string;
  createdAt: Date;
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

export interface ICommunityHotTopicsResponse {
  topics: {
    authorProfilePicture: string;
    id: number;
    lastReplyDate: Date;
    lastReplyPreview: string;
    replies: number;
    title: string;
  }[];
}

export interface ICommunityTopicsResponse {
  topics: {
    authorProfilePicture: string;
    id: number;
    lastReplyDate: Date;
    lastReplyPreview: string;
    replies: number;
    title: string;
  }[];
}
