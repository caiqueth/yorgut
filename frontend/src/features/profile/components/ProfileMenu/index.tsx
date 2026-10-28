import React, { useState } from "react";
import { ILateralButtonItem, LateralButtonList, placeholderItem } from "components/LateralButtonList";

import { Box } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useMe, useMeAmFriend } from "api/me/queries";
import { useAddFriend, useRemoveFriend, useReplyFriendship } from "api/friendship/queries";
import { useProfilePreview } from "contexts/ProfilePreviewContext";
import { useUser } from "api/user/queries";
import { ConfirmDialog } from "components/ConfirmDialog";

interface IProfileMenuProps {
  userId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ userId }) => {
  const navigate = useNavigate();
  const { data: me } = useMe();
  const { data: friendshipState, isLoading: isLoadingFriendship } = useMeAmFriend(userId);
  // The preview shows the menu as a stranger sees it.
  const { previewing } = useProfilePreview();
  const isMe = me?.id === userId && !previewing;
  const { mutate: replyFriendship, isPending: isPendingReply, variables: replying } = useReplyFriendship();
  const { mutate: removeFriend, isPending: isPendingRemoval } = useRemoveFriend(userId);
  const { mutate: addFriend, isPending: isPendingAdd } = useAddFriend();
  const { data: user } = useUser(userId);
  // "desfazer amizade" asks first.
  const [confirmingRemoval, setConfirmingRemoval] = useState(false);

  const handleDoNothing = () => { };

  const handleAddFriend = () => {
    addFriend({
      userRequesteeId: userId,
    })
  };

  const handleRemoveFriend = () => {
    if (!friendshipState?.friendshipId) return;
    removeFriend();
  };

  const handleAcceptFriendship = () => {
    if (!friendshipState?.friendshipId) return;
    replyFriendship({
      friendshipId: friendshipState.friendshipId,
      accept: true,
    })
  };

  const handleRejectFriendship = () => {
    if (!friendshipState?.friendshipId) return;
    replyFriendship({
      friendshipId: friendshipState.friendshipId,
      accept: false,
    })
  };

  const handleGoToRoute = (route: string) => {
    navigate(route);
  };

  const itemsList: ILateralButtonItem[] = [
    {
      text: "perfil", avatar: "🙂", onClick: () => handleGoToRoute(`/profile/${userId}`)
    },
    {
      text: "recados", avatar: "📝", onClick: () => handleGoToRoute(`/profile/${userId}/scraps`)
    },
    { text: "álbum", avatar: "📷", onClick: () => handleGoToRoute("albums") },
    {
      text: "escrever depoimento",
      avatar: "🖊️",
      onClick: () => handleGoToRoute("testimonials"),
    },
  ];

  if (previewing) {
    // What a stranger gets; it does nothing here.
    itemsList.unshift({ text: "adicionar", avatar: "➕" });
  } else if (!me || (!isMe && isLoadingFriendship)) {
    // Not known yet whether this is "editar perfil" or a friendship action.
    itemsList.unshift(placeholderItem);
  } else if (friendshipState) {
    if (friendshipState.amIFriend) {
      itemsList.unshift({ text: "desfazer amizade", avatar: "➖", onClick: () => setConfirmingRemoval(true), loading: isPendingRemoval })
    } else if (friendshipState.isPendingTheirReply) {
      itemsList.unshift({ text: "pedido enviado", avatar: "🕛", onClick: () => handleDoNothing() })
    } else if (friendshipState.isPendingMyReply) {
      itemsList.unshift({
        text: "aceitar amizade",
        avatar: "✅",
        onClick: () => handleAcceptFriendship(),
        loading: isPendingReply && replying?.accept,
        disabled: isPendingReply,
      })
      itemsList.unshift({
        text: "recusar amizade",
        avatar: "❌",
        onClick: () => handleRejectFriendship(),
        loading: isPendingReply && !replying?.accept,
        disabled: isPendingReply,
      })
    } else {
      itemsList.unshift({ text: "adicionar", avatar: "➕", onClick: () => handleAddFriend(), loading: isPendingAdd })
    }
  }

  // Last unshift, so it's always the first option on your own profile.
  if (isMe) {
    itemsList.unshift({ text: "editar perfil", avatar: "✏️", onClick: () => handleGoToRoute(`/profile/${userId}/edit`) })
  }

  return (
    <Box>
      <LateralButtonList items={itemsList} />
      <ConfirmDialog
        open={confirmingRemoval}
        title={"desfazer amizade"}
        text={`Tem certeza que quer desfazer a amizade com ${user?.firstName ?? "esta pessoa"}?`}
        confirmLabel={"desfazer"}
        onConfirm={handleRemoveFriend}
        onClose={() => setConfirmingRemoval(false)}
      />
    </Box>
  );
};
