import React, { useState, useEffect } from "react";
import { NavLink, useLocation } from "react-router-dom";
import Cookies from "js-cookie";
import Logo from "../../assets/Malang.png";

const Navbar = () => {
  const [isMenuOpen, setMenuOpen] = useState(false);
  const [isScrolled, setIsScrolled] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const location = useLocation();

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 50);
    };
    checkLoginStatus();

    window.addEventListener("scroll", handleScroll);

    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, []);

  const isTransparent = () => {
    const transparentRoutes = ["/"];
    return transparentRoutes.includes(location.pathname) && !isScrolled;
  };

  const toggleMenu = () => {
    setMenuOpen(!isMenuOpen);
  };

  const checkLoginStatus = () => {
    const token = Cookies.get("token");
    setIsLoggedIn(!!token);
  };

  const handleLogout = () => {
    Cookies.remove("token");
    setIsLoggedIn(false);
  };

  return (
    <nav
      className={`fixed top-0 w-full z-10 transition-colors duration-300 ${
        isTransparent() ? "bg-transparent" : "bg-emerald-700"
      }`}
    >
      <div className="max-w-screen-xl flex flex-wrap items-center justify-between mx-auto p-4">
        <a href="/" className="flex items-center space-x-3 rtl:space-x-reverse">
          <img src={Logo} className="h-11 w-11" alt="Desa Ngebruk Logo" />
          <span className="self-center text-2xl font-semibold whitespace-nowrap text-white dark:text-white">
            Desa Ngebruk
          </span>
        </a>
        <button
          data-collapse-toggle="navbar-default"
          type="button"
          className="inline-flex items-center p-2 w-10 h-10 justify-center text-sm text-gray-500 rounded-lg md:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
          aria-controls="navbar-default"
          aria-expanded={isMenuOpen ? "true" : "false"}
          onClick={toggleMenu}
        >
          <span className="sr-only">Open main menu</span>
          <svg
            className="w-5 h-5"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 17 14"
          >
            <path
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M1 1h15M1 7h15M1 13h15"
            />
          </svg>
        </button>
        <div
          className={`${
            isMenuOpen ? "block" : "hidden"
          } w-full md:block md:w-auto`}
          id="navbar-default"
        >
          <div
            className={`${
              isMenuOpen ? "block" : "hidden"
            } w-full md:block md:w-auto font-medium flex flex-col p-4 md:p-0 mt-4 border border-gray-100 rounded-lg bg-emerald-600 md:flex-row md:space-x-8 rtl:space-x-reverse md:mt-0 md:border-0 md:bg-transparent dark:bg-gray-800 md:dark:bg-gray-900 dark:border-gray-700`}
          >
            <ul className="flex flex-col mt-4 md:flex-row md:space-x-8 md:mt-0">
              <li>
                <NavLink
                  to="/"
                  className={({ isActive }) =>
                    isActive
                      ? "text-green-300"
                      : "text-white hover:text-green-400"
                  }
                >
                  Home
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/profile-desa"
                  className={({ isActive }) =>
                    isActive
                      ? "text-green-300"
                      : "text-white hover:text-green-400"
                  }
                >
                  Profile Desa
                </NavLink>
              </li>
              {isLoggedIn ? (
                <>
                  <li>
                    <NavLink
                      to="/blogger/dashboard"
                      className={({ isActive }) =>
                        isActive
                          ? "text-green-300"
                          : "text-white hover:text-green-400"
                      }
                    >
                      Dashboard Admin
                    </NavLink>
                  </li>
                  <li>
                    <button
                      onClick={handleLogout}
                      className="text-white hover:text-green-400"
                    >
                      Logout
                    </button>
                  </li>
                </>
              ) : (
                <li>
                  <NavLink
                    to="/blogger/loginakun"
                    className={({ isActive }) =>
                      isActive
                        ? "text-green-300"
                        : "text-white hover:text-green-400"
                    }
                  >
                    Login
                  </NavLink>
                </li>
              )}
            </ul>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
