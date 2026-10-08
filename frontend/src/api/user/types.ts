export interface IUsersResponse {
  users: {
    birthDate: string;
    createdAt: string;
    email: string;
    firstName: string;
    friendCount: number;
    id: number;
    lastName: string;
    updatedAt: string;
  }[];
}

export interface IUserResponse {
  birthDate: string;
  createdAt: string;
  email: string;
  firstName: string;
  friendCount: number;
  id: number;
  lastName: string;
  updatedAt: string;
}

export interface IUserFriendsResponse {
  friends: {
    birthDate: string;
    firstName: string;
    friendCount: number;
    id: number;
    profilePicture: string;
  }[];
}

export interface IUserMembershipsResponse {
  membershipCount: number;
  communities: {
    firstName: number;
    id: number;
    membershipCount: number;
    name: string;
    profilePicture: string;
  }[];
}

export interface IUserProfileResponse {
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

export interface ISearchUsersResponse {
  page: number;
  pageSize: number;
  count: number;
  users: {
    id: number;
    firstName: string;
    birthDate: string;
    profilePicture: string;
    friendCount: number;
  }[];
}