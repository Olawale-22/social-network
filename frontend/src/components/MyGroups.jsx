import React, { useContext, useEffect, useRef } from 'react';
import '../styles/group.css';
import { useState } from 'react';
import { useSelector } from 'react-redux';
import GroupForm from './GroupForm';
import GroupSearch from './GroupSearch';
import GroupCard from './GroupCard';

const MyGroups = () => {
	const groups = useSelector(state => state.data.groups);
	const [checkedUsers, setCheckedUsers] = useState([]);
	const [searchInput, setSearchInput] = useState("");
	const [results, setResults] = useState([]);
	const userId = localStorage.getItem("user_id");

	const handleSearch = searchTerm => {
		const searchedGroups = groups.filter(group => searchTerm !== "" && group.Name.toLowerCase().includes(searchTerm.toLowerCase()));
		setResults(searchedGroups);
	}

	const isGroupMember = checkedGroup => {
		return checkedGroup.Mentions.includes(parseInt(userId)) || checkedGroup.Admin_id === parseInt(userId);
	}


	return (
		<>
			<div className="search-side">
				<div className="search">
					<GroupSearch searchInput={searchInput} setSearchInput={setSearchInput} onSearch={handleSearch} />
				</div>

				<div className="form-group">
					<h2>Create group</h2>
					<GroupForm checkedUsers={checkedUsers} setCheckedUsers={setCheckedUsers} />
				</div>
			</div>

			<div className="groups">
				{groups && groups.length > 0 ? (
					searchInput === "" ? (
						groups.map(group => (
							<GroupCard key={group.ID} userId={userId} group={group} isMember={isGroupMember(group)} />
						))
					) : (
						results.length > 0 ? (
							results.map(item => (
								<GroupCard key={item.ID} userId={userId} group={item} isMember={isGroupMember(item)} />
							))
						) : (
							<span>No results found</span>
						)
					)
				) : (
					<span>No group data</span>
				)}
			</div>
		</>
	);
}

export default MyGroups;