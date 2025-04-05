export interface IUser {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  birthDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface IFriend {
  id: string;
  firstName: string;
  birthDate: Date;
  profilePicture: string;
  frienshipsCount: number;
}

export interface ICommunity {
  id: number;
  name: string;
  picture: string;
  members: number;
}
