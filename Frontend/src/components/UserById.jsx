import React, { useEffect, useState } from "react";
import axios from "axios";
import Cookies from "js-cookie";
import { useNavigate } from "react-router-dom";

const UserById = () => {
  const [userData, setUserData] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const token = Cookies.get("token");
    const fetchUserData = async () => {
      try {
        const response = await axios.get(
          "https://www.ngebruk.com/api-blog-ngebruk/user-id",
          {
            headers: {
              Authorization: `${token}`,
            },
          }
        );
        setUserData(response.data);
      } catch (err) {
        setError(err);
        console.error("Error fetching data: ", err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserData();
  }, []);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  const handleEdit = (id_user) => {
    navigate(`/blogger/userdetail/updateuser/${id_user}`);
  };

  return (
    <div className="">
      {userData ? (
        <div className="flex justify-center mt-32">
          <div className="w-full max-h-screen max-w-md bg-gray-150 border shadow-lg border-gray-200 rounded-xl shadow-slate-700 dark:bg-gray-800 dark:border-gray-700">
            <div className="flex flex-col item-start p-4 pt-10 pb-20 ml-9">
              <h5 className="mb-5 text-4xl font-sans font-bold text-black dark:text-white">
                Your Account
              </h5>
              <div className="flex mb-1">
                <span className="text-xl font-serif text-black dark:text-white w-28">
                  Nama
                </span>
                <span className="text-xl font-serif text-black dark:text-white">
                  {userData.nama}
                </span>
              </div>
              <div className="flex mb-1">
                <span className="text-xl font-serif text-black dark:text-white w-28">
                  Username
                </span>
                <span className="text-xl font-serif text-black dark:text-white">
                  {userData.username}
                </span>
              </div>
              <div className="flex">
                <span className="text-xl font-serif text-black dark:text-white w-28">
                  E-mail
                </span>
                <span className="text-xl font-serif text-black dark:text-white">
                  {userData.email}
                </span>
              </div>
              <div className="mt-2">
                {" "}
                <span className=" font-serif font-semibold bg-green-100 text-green-800 text-xlfont-medium me-2 px-2.5 py-0.5 rounded dark:bg-gray-700 dark:text-green-400 border border-green-400">
                  {userData.role}
                </span>
              </div>

              <div className="mt-4">
                <button
                  onClick={() => handleEdit(userData.id_user)}
                  type="button"
                  className=" px-4 py-2  text-white bg-yellow-300 border border-yellow-300 focus:outline-none hover:bg-yellow-200 focus:ring-4 focus:ring-yellow-200 font-medium rounded-lg text-sm  me-2 mb-2 dark:bg-yellow-800 dark:text-yellow dark:border-yellow-600 dark:hover:bg-yellow-700 dark:hover:border-yellow-600 dark:focus:ring-yellow-700"
                >
                  Edit
                </button>
                <button
                  onClick={() =>
                    navigate("/blogger/userdetail/change-password")
                  }
                  className="inline-flex items-center px-4 py-2 text-sm font-medium text-center text-white bg-red-700 rounded-lg hover:bg-red-800 focus:ring-4 focus:outline-none focus:ring-red-300 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-800"
                >
                  Ganti Password
                </button>
              </div>
            </div>
          </div>
        </div>
      ) : (
        <p>Data pengguna tidak ditemukan.</p>
      )}
    </div>
  );
};

export default UserById;
