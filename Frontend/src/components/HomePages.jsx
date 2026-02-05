import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import axios from "axios";
import { format } from "date-fns";
import ngebruk from "./img/ngebruk.jpg";
const HomePage = () => {
  const [blogs, setBlogs] = useState([]);

  const [startDate, setStartDate] = useState("");

  const [endDate, setEndDate] = useState("");

  const [searchQuery, setSearchQuery] = useState("");

  const [selectedCategory, setSelectedCategory] = useState("");

  const [currentPage, setCurrentPage] = useState(1);

  const [totalPages, setTotalPages] = useState(0);

  const pageSize = 6;

  const apiUrl = "http://localhost:8080/api-blog-ngebruk";

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(
          `${apiUrl}/blogger?page=${currentPage}&pageSize=${pageSize}&search=${encodeURIComponent(
            searchQuery
          )}&startDate=${encodeURIComponent(
            startDate
          )}&endDate=${encodeURIComponent(
            endDate
          )}&category=${encodeURIComponent(selectedCategory)}`
        );
        // / Urutkan data dari yang terbaru
        const sortedData = response.data.bloggers.sort((a, b) => {
          return new Date(b.created_at) - new Date(a.created_at);
        });

        setBlogs(sortedData);
        setTotalPages(response.data.totalPages);

        if (response.data.bloggers.length === 0) {
          setTotalPages(0);
        } else {
          setTotalPages(response.data.totalPages);
        }
      } catch (error) {
        console.error("Terjadi kesalahan!", error);
      }
    };

    fetchData();
  }, [
    currentPage,
    pageSize,
    searchQuery,
    startDate,
    endDate,
    selectedCategory,
    apiUrl,
  ]);

  useEffect(() => {
    document.title = "Home Ngebruk ";
    const fetchCategories = async () => {
      try {
        const response = await axios.get(`${apiUrl}/kategori`);
        // Menyimpan daftar kategori ke dalam state
        setCategories(response.data);
      } catch (error) {
        console.error("Terjadi kesalahan saat mengambil kategori!", error);
      }
    };

    fetchCategories();
  }, [apiUrl]);

  const [categories, setCategories] = useState([]);

  const handleCategoryChange = (e) => {
    setSelectedCategory(e.target.value);
    setCurrentPage(1);
  };

  const handlePrevious = () => {
    setCurrentPage((current) => Math.max(current - 1, 1));
  };

  const handleNext = () => {
    setCurrentPage((current) => Math.min(current + 1, totalPages));
  };

  const sliceword = (description) => {
    const words = description.split(" ");
    const first10Words = words.slice(0, 9).join(" ");

    return first10Words;
  };

  const formatTanggal = (created_at) => {
    return format(new Date(created_at), " dd MMMM yyyy HH:mm:ss");
  };

  const stripHtml = (html) => {
    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = html;
    return tempDiv.textContent || tempDiv.innerText || "";
  };

  return (
    <>
      <section
        className="bg-gray-500 bg-blend-multiply"
        style={{
          backgroundImage: `url(${ngebruk})`,
          backgroundSize: "cover",
          backgroundPosition: "center",
        }}
      >
        <div className="px-4 mx-auto max-w-screen-xl text-center py-24 lg:py-56">
          <h1 className="mb-4 text-4xl font-extrabold tracking-tight leading-none text-white md:text-5xl lg:text-6xl">
            Selamat Datang Di Desa Ngebruk
          </h1>
          <p className="mb-8 text-lg font-normal text-gray-300 lg:text-xl sm:px-16 lg:px-48">
            Semua berita di desa ngebruk akan tersampaikan disini
          </p>
        </div>
      </section>

      <div className="md:p-8 p-2">
        <div className="flex flex-col md:flex-row items-start justify-between">
          <div className="flex flex-col md:flex-row items-start space-x-4 md:space-x-2">
            <div className="p-2 md:p-0 flex-grow mr-2">
              <label
                htmlFor="startDate"
                className="block text-sm font-medium text-gray-700"
              >
                Mulai Tanggal
              </label>
              <input
                id="startDate"
                type="date"
                className="rounded-lg w-full"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
              />
            </div>

            <div className="flex-grow">
              <label
                htmlFor="endDate"
                className="block text-sm font-medium text-gray-700"
              >
                Sampai Tanggal
              </label>
              <input
                id="endDate"
                type="date"
                className="rounded-lg w-full -ml-2"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
              />
            </div>
          </div>

          <div className="flex flex-col md:flex-row items-start md:items-end space-x-4 mt-4 lg:mt-0">
            {/* Bagian Pencarian dan Kategori */}
            <input
              type="text"
              placeholder="Cari Blogger"
              value={searchQuery}
              onChange={(e) => {
                setSearchQuery(e.target.value);
                setCurrentPage(1); // Reset ke halaman pertama setiap kali ada pencarian baru
              }}
              className="p-2 ml-2 border rounded w-full md:w-auto"
            />

            <select
              value={selectedCategory}
              onChange={handleCategoryChange}
              className=" border rounded w-full md:w-auto border-gray-500 mt-2"
            >
              <option value="">Pilih Kategori</option>
              {categories.map((category) => (
                <option key={category.id_kategori} value={category.id_kategori}>
                  {category.kategori_name}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div className="text-center text-4xl font-serif">Berita Terkini</div>
        <div className="flex justify-center">
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 w-full max-w-6xl">
            {blogs.map((hasil) => (
              <div key={hasil.id_blogger} className="mt-5 m-1 md:m-1">
                <Link
                  to={`/blogger_detail/${hasil.id_blogger}`}
                  className="transition ease-in-out delay-150 hover:-translate-y-1 flex flex-col items-start bg-white border border-gray-200  shadow-slate-800 shadow-lg hover:bg-gray-300 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-700"
                >
                  <img
                    className="object-cover w-full h-64 "
                    src={hasil.heading_bloger}
                    alt=""
                  />
                  <div className="flex flex-col justify-start p-4 sm:h-64 lg:h-64 md:h-80 leading-normal">
                    <h5 className="text-xl font-bold tracking-tight text-emerald-500 dark:text-white">
                      {hasil.name_blog}
                    </h5>
                    <p className=" text-md font-medium">
                      Penulis: {hasil.user.nama}
                    </p>
                    <p className=" text-sm font-medium text-gray-500">
                      {formatTanggal(hasil.created_at)}
                    </p>
                    <p className=" text-sm font-bold text-green-600">
                      {hasil.kategori.kategori_name}
                    </p>
                    <p className="mb-3 font-2xl text-gray-700 dark:text-gray-400">
                      {stripHtml(sliceword(hasil.fill_blogger))}
                      <span className=" font-bold hover:text-blue-800 text-blue-500">
                        {" "}
                        . . .
                      </span>
                    </p>
                  </div>
                </Link>
              </div>
            ))}
          </div>
        </div>
        <div className="flex justify-center mt-4 mb-5">
          <button
            onClick={handlePrevious}
            disabled={currentPage === 1}
            className="flex items-center justify-center px-2 sm:px-4 h-8 sm:h-10 me-2 sm:me-3 text-sm sm:text-base font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
          >
            <svg
              className="w-3.5 h-3.5 me-2 rtl:rotate-180"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 14 10"
            >
              <path
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 5H1m0 0 4 4M1 5l4-4"
              />
            </svg>
            Previous
          </button>
          <span className=" text-xs sm:text-sm font-serif m-2 mr-3 sm:mr-5">
            Page {currentPage} / {totalPages}
          </span>

          <button
            onClick={handleNext}
            disabled={
              currentPage === totalPages ||
              (blogs.length < pageSize && currentPage === totalPages - 1)
            }
            className="flex items-center justify-center px-2 sm:px-4 h-8 sm:h-10 text-sm sm:text-base font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
          >
            Next
            <svg
              className="w-3.5 h-3.5 ms-2 rtl:rotate-180"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 14 10"
            >
              <path
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M1 5h12m0 0L9 1m4 4L9 9"
              />
            </svg>
          </button>
        </div>
      </div>
    </>
  );
};

export default HomePage;
