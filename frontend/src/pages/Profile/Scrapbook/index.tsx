import { Scrapbook } from "components/organisms/Scrapbook";
import React from "react";
import { useParams } from "react-router-dom";

interface IProps { }
export const UserScrapbook: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  return (
    <Scrapbook userId={userId}></Scrapbook>
  )
};
