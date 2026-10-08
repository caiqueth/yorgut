export interface ICommunityResponse {
  name: string;
  category: string;
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

export interface ICreateCommunityRequest {
  name: string;
  category: string;
  // Optional community picture (PNG or JPEG, up to 1 MB and 1024x1024).
  picture: File | null;
}

export interface IUpdateCommunityPayload {
  name: string;
  category: string;
  // An empty description clears it.
  description: string;
}

export interface ICreateCommunityResponse {
  data: { id: number };
}

export interface ISearchCommunitiesResponse {
  page: number;
  pageSize: number;
  count: number;
  communities: {
    id: number;
    name: string;
    category: string;
    profilePicture: string | null;
    memberCount: number;
  }[];
}
