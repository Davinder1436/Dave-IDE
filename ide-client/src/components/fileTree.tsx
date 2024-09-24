import React, { useEffect } from 'react';
import { useFileTree } from "./../Context/FileTreeContext"; // Adjust the import path

const FileTree: React.FC = () => {
    const { fileTree, setCurrentFile, updateFileTree } = useFileTree();

    const handleFileClick = (fileName: string) => {
        setCurrentFile(fileName);
        // Additional logic for when a file is clicked can be added here
    };

    const renderFileTree = (nodes: typeof fileTree) => {

        useEffect(() => {
            
        })
        return nodes.map((node) => (
            <li key={node.name}>
                {node.type === 'directory' ? (
                    <>
                        <span>{node.name}</span>
                        <ul>{renderFileTree(node.children || [])}</ul>
                    </>
                ) : (
                    <span onClick={() => handleFileClick(node.name)}>{node.name}</span>
                )}
            </li>
        ));

    };

    return (
        <div>
            <h1>File Tree</h1>
            <ul>{renderFileTree(fileTree)}</ul>
            <button onClick={updateFileTree}>Refresh File Tree</button>
        </div>
    );
};

export default FileTree;
