import { UpdateProfileRequestUpdateFieldsEnum, type UserInput } from "@/api/generated/models";

type ProfileFields = {
  nickname: string;
  email: string;
};

type ProfileUpdate = {
  user: UserInput;
  updateFields: UpdateProfileRequestUpdateFieldsEnum[];
};

function buildProfileUpdate(saved: ProfileFields, current: ProfileFields): ProfileUpdate | null {
  const user: UserInput = {};
  const updateFields: UpdateProfileRequestUpdateFieldsEnum[] = [];

  if (current.nickname !== saved.nickname) {
    user.nickname = current.nickname;
    updateFields.push(UpdateProfileRequestUpdateFieldsEnum.Nickname);
  }
  if (current.email !== saved.email) {
    user.email = current.email;
    updateFields.push(UpdateProfileRequestUpdateFieldsEnum.Email);
  }

  if (updateFields.length === 0) {
    return null;
  }
  return { user, updateFields };
}

export { buildProfileUpdate };
export type { ProfileFields, ProfileUpdate };
