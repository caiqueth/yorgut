export interface ITopicRepliesResponse {
  topic: {
    id: number;
    authorId: number;
    authorFirstName: string;
    authorProfilePicture: string;
    title: string;
    text: string;
    createdAt: Date;
    replies: number;
  };
  replies: {
    id: string;
    authorId: number;
    authorFirstName: string;
    authorProfilePicture: string;
    title: string;
    text: string;
    createdAt: Date;
  }[];
}

export interface ITopicCreateRequest {
  communityId: number;
  title: string;
  text: string;
}

// No communityId: the backend posts the reply in the parent topic's community.
export interface ITopicReplyRequest {
  title: string;
  text: string;
}
