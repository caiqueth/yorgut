export interface ITestimonial {
  id: number;
  content: string;
  authorId: number;
  authorName: string;
  authorProfilePicture: string;
  approved: boolean;
  createdAt: Date;
}

export interface ITestimonialsResponse {
  // Paging of the approved testimonials; count is their total.
  page: number;
  pageSize: number;
  count: number;
  testimonials: ITestimonial[];
  // Only filled when the logged user is looking at their own testimonials.
  pending: ITestimonial[];
  canWrite: boolean;
}

export interface ICreateTestimonialPayload {
  content: string;
  userRecipientId: number;
}
