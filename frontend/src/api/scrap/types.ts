export interface IScrapsResponse {
  scraps: {
    id: number;
    authorId: number;
    authorName: string;
    authorProfilePicture: string;
    content: string;
    createdAt: Date;
  }[];
}

export interface ICreateScrapPayload {
  content: string;
  userRecipientId: number;
}
