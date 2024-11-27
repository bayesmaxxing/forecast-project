import { useState, useEffect, useRef } from 'react';
import { 
  Paper,
  Box,
  TextField,
  IconButton,
  Typography,
  CircularProgress,
  Alert,
  Fab,
  Zoom,
  Badge,
} from '@mui/material';
import ChatIcon from '@mui/icons-material/Chat';
import SendIcon from '@mui/icons-material/Send';
import CloseIcon from '@mui/icons-material/Close';

const AgentChat = ({
  height = '500px',
  width = '350px',
  websocketUrl = 'ws://localhost:8000/ws',
  sx = {}
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [unreadCount, setUnreadCount] = useState(0);
  const wsRef = useRef(null);
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    if(!isOpen && messages.length > 0) {
      setUnreadCount(prev => prev + 1);
    }
  }, [messages, isOpen]);

  useEffect(() => {
    if (isOpen) {
      setUnreadCount(0);
    }
  }, [isOpen]);

  useEffect(() => {
    wsRef.current = new WebSocket(websocketUrl);

    wsRef.current.onopen = () => {
      setIsConnected(true);
      setError('');
    };

    wsRef.current.onclose = () => {
      setIsConnected(false);
      setError('Disconnected from server. Please refresh the page.');
    };

    wsRef.current.onerror = () => {
      setError('Connection error occurred. Please check your connection.');
    };

    let currentResponse = '';

    wsRef.current.onmessage = (event) => {
      if (event.data === '[DONE]') {
        setMessages(prev => {
          const newMessages = [...prev];
          if (newMessages.length > 0) {
            newMessages[newMessages.length - 1] = {
              ...newMessages[newMessages.length - 1],
              content: currentResponse,
              complete: true
            };
          }
          return newMessages;
        });
        currentResponse = '';
        setIsLoading(false);
      } else {
        currentResponse += event.data;
        setMessages(prev => {
          const newMessages = [...prev];
          if (newMessages.length > 0 && !newMessages[newMessages.length - 1].complete) {
            newMessages[newMessages.length - 1] = {
              ...newMessages[newMessages.length - 1],
              content: currentResponse
            };
          }
          return newMessages;
        });
      }
    };

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [websocketUrl]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!input.trim() || !isConnected || isLoading) return;

    setMessages(prev => [
      ...prev, 
      { role: 'user', content: input, complete: true },
      { role: 'assistant', content: '', complete: false }
    ]);
    
    setIsLoading(true);
    wsRef.current.send(input);
    setInput('');
  };

  return (
    <Box sx={{ position: 'fixed', bottom: 24, right: 24, zIndex: 1000 }}>
      <Zoom in={isOpen}>
        <Paper
          elevation={4}
          sx={{
            position: 'absolute',
            bottom: '80px',
            right: 0,
            height,
            width,
            display: 'flex',
            flexDirection: 'column',
            overflow: 'hidden',
            borderRadius: 2,
            ...sx
          }}
        >
          <Box
            sx={{
              p: 2,
              bgcolor: 'primary.main',
              color: 'primary.contrastText',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}
          >
            <Typography variant="h6">Chat</Typography>
            <IconButton 
              size="small" 
              onClick={() => setIsOpen(false)}
              sx={{ color: 'inherit' }}
            >
              <CloseIcon />
            </IconButton>
          </Box>

          <Box
            sx={{
              flexGrow: 1,
              overflowY: 'auto',
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              gap: 1,
              '&::-webkit-scrollbar': {
                width: '0.4em'
              },
              '&::-webkit-scrollbar-track': {
                bgcolor: 'grey.100'
              },
              '&::-webkit-scrollbar-thumb': {
                bgcolor: 'grey.400',
                borderRadius: '4px'
              }
            }}
          >
            {messages.map((message, index) => (
              <Box
                key={index}
                sx={{
                  maxWidth: '80%',
                  alignSelf: message.role === 'user' ? 'flex-end' : 'flex-start',
                  bgcolor: message.role === 'user' ? 'primary.main' : 'primary.secondary',
                  color: message.role === 'user' ? 'primary.contrastText' : 'text.primary',
                  p: 2,
                  borderRadius: 2,
                }}
              >
                <Typography>
                  {message.content || (isLoading && index === messages.length - 1 ? '...' : '')}
                </Typography>
              </Box>
            ))}
            <div ref={messagesEndRef} />
          </Box>

          {error && (
            <Alert severity="error" sx={{ mx: 2, mb: 2 }}>
              {error}
            </Alert>
          )}

          <Box
            component="form"
            onSubmit={handleSubmit}
            sx={{
              p: 2,
              borderTop: 1,
              borderColor: 'divider',
            }}
          >
            <Box sx={{ display: 'flex', gap: 1 }}>
              <TextField
                fullWidth
                size="small"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Ask a question..."
                disabled={!isConnected || isLoading}
              />
              <IconButton
                type="submit"
                color="primary"
                disabled={!isConnected || isLoading || !input.trim()}
              >
                {isLoading ? <CircularProgress size={24} /> : <SendIcon />}
              </IconButton>
            </Box>
          </Box>
        </Paper>
      </Zoom>

      <Badge
        badgeContent={unreadCount}
        color="error"
        sx={{
          '& .MuiBadge-badge': {
            right: 8,
            top: 8,
          }
        }}
      >
        <Fab 
          color="primary"
          onClick={() => setIsOpen(!isOpen)}
          sx={{ boxShadow: 3 }}
        >
          <ChatIcon />
        </Fab>
      </Badge>
    </Box>
  );
};

export default AgentChat;
