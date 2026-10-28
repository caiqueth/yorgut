import { Home } from "features/home/pages/Home";
import { Profile } from "features/profile/pages/Profile";

import { ProtectedRoute } from "routes/ProtectedRoute";
import { Outlet, RouteObject } from "react-router-dom";
import { Community } from "features/community/pages/Community";
import { NotFound } from "features/errors/pages/NotFound";
import { BadBadServer } from "features/errors/pages/BadBadServer";
import { CommunityMembers } from "features/community/pages/Members";
import { JoinRequests } from "features/community/pages/JoinRequests";
import { InviteFriends } from "features/community/pages/Invite";
import { Polls } from "features/community/pages/Polls";
import { Poll } from "features/community/pages/Poll";
import { NewPoll } from "features/community/pages/NewPoll";
import { EditCommunity } from "features/community/pages/Edit";
import { UserFriendsList } from "features/profile/pages/Friends";
import { UserFansList } from "features/profile/pages/Fans";
import { UserCommunitiesList } from "features/profile/pages/Communities";
import { UserPhotoAlbumsList } from "features/profile/pages/PhotoAlbums";
import { UserScrapbook } from "features/profile/pages/Scrapbook";
import { UserTestimonialsList } from "features/profile/pages/Testimonials";
import { EditProfile } from "features/profile/pages/Edit";
import { Friends } from "features/friends/pages/Friends";
import { Communities } from "features/communities/pages/Communities";
import { Search } from "features/search/pages/Search";
import { MyAccount } from "features/account/pages/MyAccount";
import { DeleteAccount } from "features/account/pages/DeleteAccount";
import { Help } from "features/info/pages/Help";
import { HelpAbout } from "features/info/pages/HelpAbout";
import { HelpRules } from "features/info/pages/HelpRules";
import { HelpTerms } from "features/info/pages/HelpTerms";
import { About } from "features/info/pages/About";
import { WhatsNew } from "features/info/pages/WhatsNew";
import { Media } from "features/info/pages/Media";
import { MyScrapbook } from "features/scrapbook/pages/MyScrapbook";
import { Login } from "features/auth/pages/Login";
import { MainPage as UserProfileMainPage } from "features/profile/pages/MainPage";
import { MainPage as CommunityMainPage } from "features/community/pages/MainPage";
import { Forum } from "features/community/pages/Forum";
import { Topic } from "features/community/pages/Topic";
import { NewTopic } from "features/community/pages/NewTopic";
import { NewCommunity } from "features/communities/pages/NewCommunity";
import { MyMemberships } from "features/communities/pages/MyMemberships";
import { MyCommunities } from "features/communities/pages/MyCommunities";
import { DiscoverCommunities } from "features/communities/pages/Discover";

// Every page but login, about and the error pages needs a session: they all
// sit under one ProtectedRoute layout.
const protectedRoutes: RouteObject[] = [
  { path: "/", element: <Home /> },
  { path: "/friends", element: <Friends /> },
  {
    path: "/communities",
    element: <Communities />,
    children: [
      { index: true, element: <MyMemberships /> },
      { path: "new", element: <NewCommunity /> },
      { path: "owned", element: <MyCommunities /> },
      { path: "discover", element: <DiscoverCommunities /> },
    ],
  },
  { path: "/scraps", element: <MyScrapbook /> },
  { path: "/search", element: <Search /> },
  { path: "/my-account", element: <MyAccount /> },
  { path: "/my-account/delete", element: <DeleteAccount /> },
  {
    path: "/help",
    element: <Help />,
    children: [
      { index: true, element: <HelpAbout /> },
      { path: "rules", element: <HelpRules /> },
      { path: "terms", element: <HelpTerms /> },
    ],
  },
  { path: "/whats-new", element: <WhatsNew /> },
  { path: "/media", element: <Media /> },
  {
    path: "/profile/:userId",
    element: <Profile />,
    children: [
      { index: true, element: <UserProfileMainPage /> },
      { path: "scraps", element: <UserScrapbook /> },
      { path: "friends", element: <UserFriendsList /> },
      { path: "fans", element: <UserFansList /> },
      { path: "communities", element: <UserCommunitiesList /> },
      { path: "testimonials", element: <UserTestimonialsList /> },
      { path: "edit", element: <EditProfile /> },
      { path: "albums", element: <UserPhotoAlbumsList /> },
    ],
  },
  {
    path: "/community/:communityId",
    element: <Community />,
    children: [
      { index: true, element: <CommunityMainPage /> },
      { path: "members", element: <CommunityMembers /> },
      { path: "requests", element: <JoinRequests /> },
      { path: "invite", element: <InviteFriends /> },
      { path: "edit", element: <EditCommunity /> },
      { path: "forum", element: <Forum /> },
      { path: "forum/topic/:topicId", element: <Topic /> },
      { path: "forum/new-topic", element: <NewTopic /> },
      { path: "polls", element: <Polls /> },
      { path: "polls/new", element: <NewPoll /> },
      { path: "polls/:pollId", element: <Poll /> },
    ],
  },
];

export const Routes: RouteObject[] = [
  {
    element: (
      <ProtectedRoute>
        <Outlet />
      </ProtectedRoute>
    ),
    children: protectedRoutes,
  },
  { path: "/login", element: <Login /> },
  // Public, like /login: the footer of every page links here.
  { path: "/about", element: <About /> },
  { path: "/not-found", element: <NotFound /> },
  { path: "/bad-server", element: <BadBadServer /> },
  { path: "*", element: <NotFound /> },
];
