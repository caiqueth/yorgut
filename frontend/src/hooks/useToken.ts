import { parseJwt } from "helpers/parseJwt";

interface IToken {
  iss: string;
  azp: string;
  aud: string;
  sub: string;
  email: string;
  email_verified: boolean;
  nbf: number;
  name: string;
  picture: string;
  given_name: string;
  family_name: string;
  iat: number;
  exp: number;
  jti: string;
}

export const useToken = () => {
  const token = localStorage.getItem("token");
  try {
    if (!token) return;
    return parseJwt(token) as IToken;
  } catch (e) {
    console.error("Error parsing JWT Token =>", e);
    return;
  }
};
