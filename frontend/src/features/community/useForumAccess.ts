import { useCommunity } from "api/community/queries";
import { useMeAmMember } from "api/me/queries";

// Whether the logged user can read the community's forum: anyone in a public
// community, only members in a private one. The backend enforces the same
// rule (403); this avoids asking for what it would refuse.
export const useForumAccess = (communityId: number) => {
  const { data: community, isLoading: isLoadingCommunity } = useCommunity(communityId);
  const { data: membership, isLoading: isLoadingMembership } = useMeAmMember(communityId);

  return {
    canRead: !!community && (!community.private || !!membership?.amIMember),
    isLoading: isLoadingCommunity || isLoadingMembership,
  };
};
