import React from "react";
import { useCommunityMembers } from "api/user/queries";
import { Button } from "components/atoms/Button";
import { useNavigate } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";

interface IMembersListProps {
  communityId: number;
}
export const MembersList: React.FC<IMembersListProps> = ({ communityId }) => {
  const { data: members } = useCommunityMembers(communityId);
  const navigate = useNavigate();

  const handleGoToMembers = () => {
    navigate("members");
  };

  const items =
    members?.members.map((member) => ({
      imageSrc: `/${member.profilePicture}`,
      url: `/profile/${member.id}`,
      text: `${member.firstName} (${member.friendCount})`,
    })) ?? [];

  return (
    <BoxWithTitle title={`membros (${members?.memberCount})`}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} />
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToMembers}>ver membros</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
