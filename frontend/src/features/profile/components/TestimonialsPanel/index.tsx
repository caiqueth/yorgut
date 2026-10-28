import { Box, Stack, Typography } from "@mui/material";
import { useTestimonials } from "api/testimonial/queries";
import { useUserProfileSection } from "api/user/queries";
import { PrivateSection } from "components/PrivateSection";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";
import { Link } from "react-router-dom";

// The panel shows the newest few; "ver todos" opens the full list.
const SHOWN = 5;

interface IProps {
  userId: number;
}
export const TestimonialsPanel: React.FC<IProps> = ({ userId }) => {
  const section = useUserProfileSection(userId, "testimonials");
  const { data, isLoading: isLoadingTestimonials, isFetching } = useTestimonials(userId, 1, SHOWN, section.canSee);
  // The request waits for the profile, which says whether it may be made.
  const isLoading = section.isLoading || isLoadingTestimonials;

  const items = (data?.testimonials ?? []).map((t) => (
    <PicAndTextHorizontal
      key={t.id}
      imageSrc={t.authorProfilePicture}
      title={t.authorName}
      titleSuffix={` em ${dayjsFromDate(t.createdAt, "D [de] MMMM [de] YYYY")}:`}
      text={t.content}
      url={`/profile/${t.authorId}`}
    />
  ));

  return (
    <BoxWithTitle
      title={section.isLoading || section.canSee ? <>depoimentos <Count value={data?.count} /></> : "depoimentos"}
      showProgress={isFetching && !isLoading}
    >
      {!section.isLoading && !section.canSee
        ? <PrivateSection firstName={section.firstName} what={"os depoimentos"} />
        : isLoading || items.length > 0
        ? <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<PicAndTextHorizontalSkeleton lineWidths={[240, 200]} />} />
        : (
          <Box sx={{ p: 1 }}>
            <Typography variant={"body2"}>Nenhum depoimento ainda.</Typography>
          </Box>
        )}
      {(section.isLoading || section.canSee) && (
        <Stack direction={"row"} justifyContent={"end"} alignItems={"center"} sx={{ px: 1, pb: 1 }}>
          <Typography variant={"body2"}>
            <Link to={`/profile/${userId}/testimonials`}>ver todos</Link>
          </Typography>
        </Stack>
      )}
    </BoxWithTitle>
  );
};
