import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { Button, Modal } from "flowbite-react";
import { format } from "date-fns";
import swal from "sweetalert";

const ListBloggerByUser = () => {
  const navigate = useNavigate();
  const [fetchStatus, setFetchStatus] = useState(true);
  const [jobs, setJobs] = useState([]);
  const [openModal, setOpenModal] = useState(false);
  const [openModals, setOpenModals] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const pageSize = 5;

  useEffect(() => {
    if (fetchStatus === true) {
      const token = Cookies.get("token"); // Ambil token dari Cookies atau tempat penyimpanan yang sesuai
      console.log(token);
      if (token) {
        axios
          .get(
            `http://localhost:8080/api-blog-ngebruk/blogger-byuser?page=${currentPage}&pageSize=${pageSize}`,
            {
              headers: {
                Authorization: `${token}`,
              },
            }
          )
          .then((response) => {
            setJobs(response.data.bloggers);
            setTotalPages(response.data.totalPages);
            window.dispatchEvent(new Event("userUpdated"));
          })
          .catch((error) => console.error("Error fetching job list:", error));
      } else {
        // Handle jika token tidak tersedia
        console.error("JWT token is missing.");
      }

      setFetchStatus(false);
    }
  }, [fetchStatus, setFetchStatus]);

  const handleEdit = (id_blogger) => {
    navigate(`/blogger/updateblogger/${id_blogger}`);
  };

  const sliceword = (description) => {
    if (description == null) {
      return "Deskripsi Kosong";
    }

    const words = description.split(" ");
    const first8Words = words.slice(0, 8).join(" ");

    return first8Words;
  };

  const formatTanggal = (created_at) => {
    return format(new Date(created_at), " dd MMMM yyyy HH:mm:ss");
  };

  //delete
  const handleDelete = async (id_blogger) => {
    const token = Cookies.get("token");
    try {
      await axios.delete(
        `http://localhost:8080/api-blog-ngebruk/blogger/${id_blogger}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      swal("Success", "Berhasil menghapus data", "success");
      setFetchStatus(true);

      // Close the modal after successful deletion
      const updatedOpenModals = [...openModals];
      const index = jobs.findIndex((job) => job.id_blogger === id_blogger);
      if (index !== -1) {
        updatedOpenModals[index] = false;
        setOpenModals(updatedOpenModals);
      }
    } catch (error) {
      swal("Error!", "Gagal menghapus data", "error");
      console.error("Error deleting blogger:", error);
    }
  };

  const handleOpenModal = (index) => {
    const newOpenModals = [...openModals];
    newOpenModals[index] = true;
    setOpenModals(newOpenModals);
  };

  const handleCloseModal = (index) => {
    const newOpenModals = [...openModals];
    newOpenModals[index] = false;
    setOpenModals(newOpenModals);
  };

  const handlePrevious = () => {
    if (currentPage > 1) {
      setCurrentPage((current) => current - 1);
      setFetchStatus(true);
    }
  };

  const handleNext = () => {
    if (currentPage < totalPages) {
      setCurrentPage((current) => current + 1);
      setFetchStatus(true);
    }
  };

  const stripHtml = (html) => {
    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = html;
    return tempDiv.textContent || tempDiv.innerText || "";
  };

  return (
    <>
      <div className="">
        <div className="p-4  mt-24">
          <h1 className="text-center font-bold text-4xl">Semua Konten</h1>
          <div className="mt-2">
            <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
              <table className="table rounded-lg overflow-hidden w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                <thead className="text-xs text-gray-700 uppercase bg-blue-400 dark:bg-gray-700 dark:text-gray-400">
                  <tr>
                    <th scope="col" className="px-6 py-3 text-white">
                      No
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Judul Konten
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Cover Konten
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Isi Konten
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Kategori
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Pembuat Konten
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Tanggal Dibuat
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Action
                    </th>
                  </tr>
                </thead>

                {jobs.map((job, numb) => (
                  <tbody>
                    <tr className="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700">
                      <th
                        scope="row"
                        className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                        key={numb}
                      >
                        {(currentPage - 1) * pageSize + numb + 1}
                      </th>
                      <td className="px-6 py-4 font-bold">
                        <h1>{job.name_blog}</h1>
                      </td>
                      <td className="px-6 py-4">
                        <div className="w-24 h-24 overflow-hidden rounded-md">
                          <img
                            className="object-cover w-full h-full"
                            src={job.heading_bloger}
                            alt=""
                          />
                        </div>
                      </td>
                      <td className="px-6 py-4 pro">
                        {stripHtml(sliceword(job.fill_blogger))}
                      </td>
                      <td className="px-6 py-4">
                        {job.kategori.kategori_name}
                      </td>
                      <td className="px-6 py-4">{job.user.nama}</td>
                      <td className="px-6 py-4">
                        {formatTanggal(job.created_at)}
                      </td>
                      <td className="px-6 py-4 space-x-4 space-x-reverse ">
                        <button
                          type="button"
                          value={job.id_blogger}
                          onClick={() => handleEdit(job.id_blogger)}
                          className=" text-white bg-yellow-300 border border-yellow-300 focus:outline-none hover:bg-yellow-100 focus:ring-4 focus:ring-yellow-200 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-yellow-800 dark:text-yellow dark:border-yellow-600 dark:hover:bg-yellow-700 dark:hover:border-yellow-600 dark:focus:ring-yellow-700"
                        >
                          Edit
                        </button>
                        <Button
                          color="failure"
                          onClick={() => handleOpenModal(numb)}
                        >
                          Delete
                        </Button>
                        <Modal
                          show={openModals[numb]}
                          onClose={() => handleCloseModal(numb)}
                        >
                          <Modal.Header>Warning</Modal.Header>
                          <Modal.Body>
                            <p className="text-base leading-relaxed text-gray-500 dark:text-gray-400">
                              Apakah anda yakin mau dihapus ?
                            </p>
                          </Modal.Body>
                          <Modal.Footer>
                            <Button
                              color="failure"
                              onClick={() => handleDelete(job.id_blogger)}
                            >
                              I accept
                            </Button>
                            <Button
                              color="gray"
                              onClick={() => handleCloseModal(numb)}
                            >
                              Decline
                            </Button>
                          </Modal.Footer>
                        </Modal>
                      </td>
                    </tr>
                  </tbody>
                ))}
              </table>
            </div>
            <div className="p-4">
              <button
                onClick={() => navigate("/blogger/createblogger")}
                type="button"
                className=" focus:outline-none text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-red-700 dark:focus:ring-red-900"
              >
                Tambah Konten
              </button>
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
              disabled={currentPage === totalPages}
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
      </div>
    </>
  );
};
export default ListBloggerByUser;
