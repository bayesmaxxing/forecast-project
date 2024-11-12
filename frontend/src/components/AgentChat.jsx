import { useState, useEffect, useRef } from 'react';
import { 
  Paper,
  Box,
  TextField,
  IconButton,
  Typography,
  CircularProgress,
  Alert,
  Card,
  CardContent
} from '@mui/material';
import SendIcon from '@mui/icons-material/Send';
import { styled } from '@mui/material/styles';

const MessageContainer = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1),
  overflowY: 'auto',
  '&::-webkit-scrollbar': {
    width: '0.4em'
  },
  '&::-webkit-scrollbar-track': {
    background: theme.palette.grey[100]
  },
  '&::-webkit-scrollbar-thumb': {
    background: theme.palette.grey[400],
    borderRadius: '4px'
  }
}));

const MessageBubble = styled(Card)(({ theme, role }) => ({
  maxWidth: '80%',
  alignSelf: role === 'user' ? 'flex-end' : 'flex-start',
  backgroundColor: role === 'user' ? theme.palette.primary.main : theme.palette.grey[100],
  color: role === 'user' ? theme.palette.primary.contrastText : theme.palette.text.primary,
  marginBottom: theme.spacing(1)
}));

// Props:
// height: string (e.g., '500px', '100vh')
// width: string (e.g., '100%', '500px')
// websocketUrl: string (e.g., 'ws://localhost:8000/ws')
// className?: string - additional CSS classes
// sx?: SxProps - MUI styling overrides
export default function ForecastChat({ 
  height = '500px', 
  width = '100%', 
  websocketUrl = 'ws://localhost:8000/ws',
  className,
  sx = {}
}) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const wsRef = useRef(null);
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

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
      setError('Connection error occurred. Please check your connection and try again.');
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
    <Paper 
      elevation={2} 
      className={className}
      sx={{ 
        height, 
        width, 
        display: 'flex', 
        flexDirection: 'column',
        overflow: 'hidden',
        ...sx 
      }}
    >
      <MessageContainer 
        sx={{ 
          flexGrow: 1, 
          p: 2,
          maxHeight: `calc(${height} - 80px)` 
        }}
      >
        {messages.map((message, index) => (
          <MessageBubble key={index} role={message.role}>
            <CardContent>
              <Typography>
                {message.content || (isLoading && index === messages.length - 1 ? '...' : '')}
              </Typography>
            </CardContent>
          </MessageBubble>
        ))}
        <div ref={messagesEndRef} />
      </MessageContainer>

      {error && (
        <Alert severity="error" sx={{ mx: 2, my: 1 }}>
          {error}
        </Alert>
      )}

      <Box component="form" onSubmit={handleSubmit} sx={{ p: 2, borderTop: 1, borderColor: 'divider' }}>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <TextField
            fullWidth
            size="small"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Ask a question..."
            disabled={!isConnected || isLoading}
            sx={{ flexGrow: 1 }}
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
  );
}
