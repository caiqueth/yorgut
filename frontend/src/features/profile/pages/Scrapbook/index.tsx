import { useIdParam } from "hooks/useIdParam";
import { Scrapbook } from "components/Scrapbook";
import React from "react";

export const UserScrapbook: React.FC = () => {
  const userId = useIdParam("userId");

  return (
    <Scrapbook userId={userId}></Scrapbook>
  )
};
