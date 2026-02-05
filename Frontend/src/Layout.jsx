import React from "react";
// ... (Impor komponen lain seperti Navbar dan Footer)

const Layout = ({ children }) => {
  return (
    <div className="flex flex-col min-h-screen">
      {/* Header atau Navbar */}
      <header>{/* Konten Header atau Navbar */}</header>

      {/* Main Content */}
      <main className="flex-grow">{children}</main>

      {/* Footer */}
      <footer className="bg-gray-200 rounded-lg shadow dark:bg-gray-900">
        {/* Konten Footer */}
      </footer>
    </div>
  );
};

export default Layout;
