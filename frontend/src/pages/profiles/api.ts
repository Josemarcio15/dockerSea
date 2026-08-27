import * as ProfileBinding from "../../../bindings/go-walis/internal/profiles/service.js";

export const listProfiles = () => ProfileBinding.ListProfiles();
export const getActiveProfile = () => ProfileBinding.GetActiveProfile();
export const saveProfile = (profile: any) =>
  ProfileBinding.SaveProfile(profile);
export const deleteProfile = (id: string) => ProfileBinding.DeleteProfile(id);
export const setActiveProfile = (id: string) =>
  ProfileBinding.SetActiveProfile(id);
