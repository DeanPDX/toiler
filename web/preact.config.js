export default (config, env, helpers) => {
  // Only set up proxy if dev server exists
  if (config.devServer) {
    console.log('Setting dev server proxy settings');
    config.devServer.proxy = [
      {
        path: '/api/**',
        target: 'http://localhost:8090/',
      }
    ];
  }
}