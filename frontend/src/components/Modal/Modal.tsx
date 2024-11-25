import React from "react";

interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    title?: string;
    children: React.ReactNode;
}

const Modal: React.FC<ModalProps> = ({ isOpen, onClose, title, children }) => {
    if (!isOpen) return null;

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50"
            onClick={onClose}
        >
            <div
                className="bg-white rounded-lg shadow-lg w-1/3 p-5"
                onClick={(e) => e.stopPropagation()}
            >
                {/* Modal Header */}
                <div className="flex justify-between items-center mb-4">
                    <h2 className="text-lg font-bold">{title}</h2>
                    <button
                        className="text-gray-500 hover:text-black"
                        onClick={onClose}
                    >
                        ✖
                    </button>
                </div>

                {/* Modal Content */}
                <div>{children}</div>
            </div>
        </div>
    );
};

export default Modal;
