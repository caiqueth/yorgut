export interface ILoginWithGooglePayload {
  // The Google ID token returned by Sign in with Google.
  credential: string;
}

export interface ILoginResponse {
  // The Yorgut session token, sent as "Authorization: Bearer <token>".
  token: string;
}
