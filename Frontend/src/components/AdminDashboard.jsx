import { PiFolderNotchPlusLight } from "react-icons/pi";
import { useEffect, useState } from "react";
import Cookies from "js-cookie";

const AdminDashboard = () => {
  const [userData, setUserData] = useState([]);

  useEffect(() => {
    document.title = "Dashboard ";

    const userDataFromCookies = Cookies.get("user");
    window.dispatchEvent(new Event("userUpdated"));

    if (userDataFromCookies) {
      // Konversi dari string ke objek
      const parsedUserData = JSON.parse(userDataFromCookies);

      setUserData(parsedUserData);

      console.log(parsedUserData);
    }
  }, []);
  return (
    <>
      <div className="p-4  mt-15 bg-gray-100  w-screen h-screen flex flex-col items-center justify-center">
        {userData ? (
          <div className="p-4 ">
            <h1 className="text-center font-bold text-2xl text-gray-800">
              Selamat Datang {userData.name}
            </h1>
            <h3 className="text-center font-semibold text-gray-500 text-xl">
              Kamu dapat melakukan Tambah Data, Edit Data dan Hapus Data
            </h3>
            <div className="flex flex-col items-center justify-center">
              <PiFolderNotchPlusLight size="64" color="gray" />
            </div>
          </div>
        ) : (
          <p>Data pengguna tidak ditemukan.</p>
        )}
      </div>
    </>
  );
};
export default AdminDashboard;
