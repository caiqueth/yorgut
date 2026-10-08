export interface IUserReview {
  cool: number;
  sexy: number;
  trustworthy: number;
}

export interface IUserReviewSummaryResponse extends IUserReview {
  reviewCount: number;
  canReview: boolean;
  myReview: IUserReview | null;
}

export type IUpsertUserReviewPayload = IUserReview;
