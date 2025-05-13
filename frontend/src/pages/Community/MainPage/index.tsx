import { Grid2 as Grid, Stack } from "@mui/material";
import React from "react";
import { useParams } from "react-router-dom";
import { ProfileHeader } from "../components/ProfileHeader";
import { ProfileDetails } from "../components/ProfileDetails";
import { MembersList } from "components/organisms/MembersList";
import { RelatedCommunitiesList } from "../components/RelatedCommunitiesList";
import { HotTopics } from "../components/HotTopics";

interface IProps { };
export const MainPage: React.FC<IProps> = () => {
  const params = useParams()
  const communityId = parseInt(params.communityId ?? "");

  return <Grid container spacing={2}>
    <Grid size={7}>
      <Stack direction={"column"} gap={2}>
        <ProfileHeader communityId={communityId} />
        <ProfileDetails communityId={communityId} />
        <HotTopics communityId={communityId} />
      </Stack>
    </Grid>
    <Grid size={5}>
      <Stack direction={"column"} gap={2}>
        <MembersList communityId={communityId} />
        <RelatedCommunitiesList communityId={communityId} />
      </Stack>
    </Grid>
  </Grid>
}