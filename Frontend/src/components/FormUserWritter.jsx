import React, { useState, useEffect } from "react";
import axios from "axios";
import Cookies from "js-cookie";
import Swal from "sweetalert2";
import { useNavigate, useParams } from "react-router-dom";

const FormUserWritter = () => {
  const [nama, setNama] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const navigate = useNavigate();
  const { id_user } = useParams();

  useEffect(() => {
    document.title = "Form User";
    if (id_user) {
      axios
        .get(`http://localhost:8080/api-blog-ngebruk/user/${id_user}`)
        .then((response) => {
          const { nama, username, email } = response.data;
          setNama(nama);
          setUsername(username);
          setEmail(email);
        })
        .catch((error) => {
          console.error("Error fetching user data:", error);
          Swal.fire("Error", "Gagal memuat data user", "error");
        });
    }
  }, [id_user]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const token = Cookies.get("token");
    const url = id_user
      ? `http://localhost:8080/api-blog-ngebruk/users/${id_user}`
      : "http://localhost:8080/api-blog-ngebruk/register";
    const method = id_user ? "patch" : "post";

    const userData = { nama, username, email };

    try {
      const response = await axios[method](url, userData, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      // Assuming the response contains a new token and user details
      const newToken = response.data.token;
      const updatedUser = response.data.user;

      // Save the new token and user details in cookies
      Cookies.set("token", newToken, { expires: 1 });
      Cookies.set("user", JSON.stringify(updatedUser), { expires: 1 });
      window.dispatchEvent(new Event("userUpdated"));

      Swal.fire(
        "Success",
        `User berhasil ${id_user ? "diperbarui" : "ditambahkan"}`,
        "success"
      );
      navigate("/blogger/userdetail");
    } catch (error) {
      console.error("Error:", error);
      Swal.fire(
        "Error",
        `Gagal ${
          id_user ? "memperbarui" : "menambahkan"
        } username atau email sudah terpakai`,
        "error"
      );
    }
  };

  return (
    <section className="p-4 mt-20 bg-white dark:bg-gray-900">
      <div className="py-8 px-4 mx-auto max-w-2xl lg:py-16">
        <h2 className="mb-4 text-xl font-bold text-gray-900 dark:text-white">
          {id_user ? "Edit" : "Tambah"} Akun
        </h2>
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 sm:grid-cols-2 sm:gap-6">
            {/* Nama */}
            <div className="sm:col-span-2">
              <label
                htmlFor="nama"
                className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >
                Nama
              </label>
              <input
                type="text"
                id="nama"
                value={nama}
                onChange={(e) => setNama(e.target.value)}
                className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                placeholder="Nama"
                required
              />
            </div>
            {/* Email */}
            <div className="sm:col-span-2">
              <label
                htmlFor="email"
                className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >
                Email
              </label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                placeholder="email@example.com"
                required
              />
            </div>
            {/* Username */}
            <div className="sm:col-span-2">
              <label
                htmlFor="username"
                className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >
                Username
              </label>
              <input
                type="text"
                id="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                placeholder="Username"
                required
              />
            </div>
          </div>
          <button
            type="submit"
            className="inline-flex items-center px-5 py-2.5 mt-4 sm:mt-6 text-sm font-medium text-center text-white bg-blue-700 rounded-lg focus:ring-4 focus:ring-blue-200 dark:focus:ring-blue-900 hover:bg-blue-800"
          >
            Update Akun
          </button>
        </form>
      </div>
    </section>
  );
};

export default FormUserWritter;
