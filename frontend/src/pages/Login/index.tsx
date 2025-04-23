import { GoogleLogin } from "@react-oauth/google";
import { useNavigate } from "react-router-dom";

export const Login = () => {
  const navigate = useNavigate();
  return (
    <>
      <GoogleLogin
        onSuccess={(credentialResponse) => {
          localStorage.setItem(
            "token",
            JSON.stringify(credentialResponse.credential),
          );
          navigate("/");
        }}
        onError={() => {
          console.log("Login Failed");
        }}
      />
    </>
  );
};
