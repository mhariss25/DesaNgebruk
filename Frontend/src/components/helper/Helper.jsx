import Cookies from "js-cookie";

export const getToken = () => {
  return Cookies.get("token");
};

export const getUserRoleFromCookie = () => {
  const userCookie = Cookies.get("user");
  if (!userCookie) return null;

  try {
    const user = JSON.parse(userCookie);
    return user.role;
  } catch (error) {
    console.error("Error parsing user data from cookie:", error);
    return null;
  }
};
