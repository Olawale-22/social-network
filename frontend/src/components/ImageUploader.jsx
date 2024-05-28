import React from 'react';

const ImageUploader = ({ onImageSelected }) => {

    const handleImageChange = (event) => {
        const image = event.target.files[0];
        onImageSelected(image);
    };

    return (
        <input
            name="image"
            id="image"
            type="file"
            accept="image/*"
            onChange={handleImageChange}
        />
    );
};

export default ImageUploader;