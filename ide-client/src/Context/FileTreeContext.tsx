import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import axios from 'axios';

// Define types for FileNode and context state
interface FileNode {
    name: string;
    type: 'file' | 'directory';
    children?: FileNode[];
}

interface FileTreeContextType {
    fileTree: FileNode[];
    currentFile: string | null;
    setCurrentFile: React.Dispatch<React.SetStateAction<string | null>>;
    updateFileTree: () => Promise<void>;
}

// Create the context with a default value
const FileTreeContext = createContext<FileTreeContextType | undefined>(undefined);

// Provider component to manage the file tree state
export const FileTreeProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [fileTree, setFileTree] = useState<FileNode[]>([]);
    const [currentFile, setCurrentFile] = useState<string | null>(null);

    // Function to fetch the file tree from the backend
    const fetchFileTree = async () => {
        try {
            const response = await axios.get<FileNode[]>('/filetree'); // Adjust endpoint as needed
            setFileTree(response.data);
        } catch (error) {
            console.error("Error fetching file tree:", error);
        }
    };

    // Update the file tree when a new file is added
    const updateFileTree = async () => {
        await fetchFileTree();
    };

    // Fetch file tree on component mount
    useEffect(() => {
        fetchFileTree();
    }, []);

    return (
        <FileTreeContext.Provider value={{ fileTree, currentFile, setCurrentFile, updateFileTree }}>
            {children}
        </FileTreeContext.Provider>
    );
};

// Custom hook to use the FileTree context
export const useFileTree = (): FileTreeContextType => {
    const context = useContext(FileTreeContext);
    if (context === undefined) {
        throw new Error("useFileTree must be used within a FileTreeProvider");
    }
    return context;
};
