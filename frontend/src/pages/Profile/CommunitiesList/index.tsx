import React from "react";

interface ICommunitiesListProps {
  userId: number;
}
export const CommunitiesList: React.FC<ICommunitiesListProps> = ({ userId }) => {
  return <>{userId}</>;
};
