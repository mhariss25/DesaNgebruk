import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { getToken, getUserRoleFromCookie } from "../helper/Helper";

const PrivateRoute = ({ children, allowedRoles }) => {
  const location = useLocation();
  const token = getToken();
  const userRole = getUserRoleFromCookie();

  if (!token) {
    // Jika tidak ada token, redirect ke halaman login
    return <Navigate to="/blogger/loginakun" state={{ from: location }} />;
  }

  if (!allowedRoles.includes(userRole)) {
    // Jika role tidak diizinkan, redirect ke halaman tidak diizinkan
    return <Navigate to="/blogger/dashboard" />;
  }

  return children;
};
export default PrivateRoute;
