import React, { useState, useEffect, type FormEvent } from "react";
import "./App.css";

interface Comment {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  content: string;
  post_id: number;
}

interface Post {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  title: string;
  content: string;
  comments: Comment[];
}

const API_URL = "http://localhost:8080"; 

function App() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [newPostTitle, setNewPostTitle] = useState("");
  const [newPostContent, setNewPostContent] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedPost, setSelectedPost] = useState<Post | null>(null); 

  const fetchPosts = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_URL}/posts/`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data: Post[] = await response.json();
      setPosts(data);
    } catch (e: any) {
      setError(e.message || "Failed to fetch posts.");
      console.error("Fetch posts error:", e);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  const handleCreatePost = async (e: FormEvent) => {
    e.preventDefault();
    if (!newPostTitle.trim() || !newPostContent.trim()) {
      setError("Title and content cannot be empty.");
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_URL}/posts/`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ title: newPostTitle, content: newPostContent }),
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(
          errorData.error || `HTTP error! status: ${response.status}`
        );
      }
      // const createdPost: Post = await response.json(); // Dữ liệu post vừa tạo
      setNewPostTitle("");
      setNewPostContent("");
      fetchPosts(); 
    } catch (e: any) {
      setError(e.message || "Failed to create post.");
      console.error("Create post error:", e);
    } finally {
      setIsLoading(false);
    }
  };

  // Hàm xử lý khi click vào một post để xem chi tiết
  const handlePostClick = (post: Post) => {
    setSelectedPost(post);
  };

  // Hàm để đóng modal/chi tiết post
  const closePostDetail = () => {
    setSelectedPost(null);
  };

  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Social Media Feed</h1>
      </header>

      <main className="app-main">
        <section className="create-post-section">
          <h2>Create New Post</h2>
          <form onSubmit={handleCreatePost} className="create-post-form">
            <div>
              <label htmlFor="postTitle">Title:</label>
              <input
                type="text"
                id="postTitle"
                value={newPostTitle}
                onChange={(e) => setNewPostTitle(e.target.value)}
                required
              />
            </div>
            <div>
              <label htmlFor="postContent">Content:</label>
              <textarea
                id="postContent"
                value={newPostContent}
                onChange={(e) => setNewPostContent(e.target.value)}
                required
              />
            </div>
            <button type="submit" disabled={isLoading}>
              {isLoading ? "Creating..." : "Create Post"}
            </button>
          </form>
        </section>

        {error && <p className="error-message">Error: {error}</p>}

        <section className="posts-section">
          <h2>Posts</h2>
          {isLoading && posts.length === 0 && <p>Loading posts...</p>}
          {!isLoading && posts.length === 0 && !error && (
            <p>No posts yet. Create one!</p>
          )}
          <div className="posts-list">
            {posts.map((post) => (
              <article
                key={post.ID}
                className="post-item"
                onClick={() => handlePostClick(post)}
              >
                <h3>{post.title}</h3>
                <p>
                  {post.content.substring(0, 100)}
                  {post.content.length > 100 ? "..." : ""}
                </p>
                <small>
                  Comments: {post.comments ? post.comments.length : 0}
                </small>
              </article>
            ))}
          </div>
        </section>
      </main>

      {/* Modal hiển thị chi tiết Post và Comments */}
      {selectedPost && (
        <div className="modal-overlay" onClick={closePostDetail}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            {" "}
            {/* Ngăn click bên trong modal đóng modal */}
            <button className="modal-close-button" onClick={closePostDetail}>
              ×
            </button>
            <h2>{selectedPost.title}</h2>
            <p className="post-full-content">{selectedPost.content}</p>
            <hr />
            <h4>
              Comments (
              {selectedPost.comments ? selectedPost.comments.length : 0}):
            </h4>
            {selectedPost.comments && selectedPost.comments.length > 0 ? (
              <ul className="comments-list">
                {selectedPost.comments.map((comment) => (
                  <li key={comment.ID} className="comment-item">
                    <p>{comment.content}</p>
                    <small>
                      Commented on:{" "}
                      {new Date(comment.CreatedAt).toLocaleString()}
                    </small>
                  </li>
                ))}
              </ul>
            ) : (
              <p>No comments yet.</p>
            )}
            {/* Bạn có thể thêm form tạo comment ở đây nếu muốn */}
          </div>
        </div>
      )}

      <footer className="app-footer">
        <p>© {new Date().getFullYear()} My Social Media App</p>
      </footer>
    </div>
  );
}

export default App;
