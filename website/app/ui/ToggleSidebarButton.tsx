"use client";
import React, { useState } from "react";
import clsx from "clsx";

const ToggleSidebarButton = () => {
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);

    const toggleSidebar = () => {
        setIsSidebarOpen(!isSidebarOpen);
    };

    return (
        <button
            onClick={toggleSidebar}
            aria-pressed={isSidebarOpen}
            aria-label="toggle sidebar"
            className="peer absolute right-0 z-10 space-y-[7px] p-4 md:hidden"
        >
            <span
                className={clsx(
                    "block h-[1px] w-7 bg-white transition-[opacity,transform]",
                    isSidebarOpen && "translate-y-[7px] -rotate-45 transform",
                )}
            ></span>
            <span
                className={clsx(
                    "block h-[1px] w-7 bg-white transition-[opacity,transform]",
                    isSidebarOpen && "rotate-45 transform",
                )}
            ></span>
            <span
                className={clsx(
                    "block h-[1px] w-7 bg-white transition-[opacity,transform]",
                    isSidebarOpen && "opacity-0",
                )}
            ></span>
        </button>
    );
};

export default ToggleSidebarButton;
