import React from 'react';

const GroupSearch = ({ searchInput, setSearchInput, onSearch }) => {
    const handleChange = event => {
        const inputValue = event.target.value;
        setSearchInput(inputValue);
        onSearch(inputValue);
    }

    return (
        <>
            <input className='search-group' type="text" placeholder='Search a group...' name="search-group" id="search-group" onChange={handleChange} value={searchInput} />
        </>
    );
}

export default GroupSearch